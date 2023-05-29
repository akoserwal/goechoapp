package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"io"
	"log"
	"net/http"
	"sync"
)

const defaultport = "3000"

type EchoResponse struct {
	AccessToken string      `json:"access_token,omitempty"`
	QueryParam  string      `json:"query_param,omitempty"`
	Body        interface{} `json:"body,omitempty"`
}

func main() {
	_ = flag.CommandLine.Parse([]string{})

	var rootCmd = &cobra.Command{
		Use:   "echoapp",
		Short: "A simple golang service with echo end-point.",
		Long:  `A simple golang service with echo end-point`,
		Run:   serve,
	}

	rootCmd.Flags().String("host", "", "host")
	rootCmd.Flags().String("port", "", "port")
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("error running command: %v", err)
	}

}
func mustGetString(flagName string, flags *pflag.FlagSet) string {
	flagVal, err := flags.GetString(flagName)
	if err != nil {
		log.Fatalf(notFoundMessage(flagName, err))
	}
	return flagVal
}
func notFoundMessage(flagName string, err error) string {
	return fmt.Sprintf("could not get flag %s from flag set: %s", flagName, err.Error())
}

func serve(cmd *cobra.Command, _ []string) {
	host := mustGetString("host", cmd.Flags())
	port := mustGetString("port", cmd.Flags())
	if port == "" {
		port = defaultport
	}
	wait := sync.WaitGroup{}
	go func() {
		err := httpServer(&wait, host, port)
		if err != nil {
			log.Println("Could not start http serving: ", err)
		}
	}()

	wait.Add(1)
	wait.Wait()
}

func httpServer(wait *sync.WaitGroup, host string, port string) error {
	defer wait.Done()
	mux := http.NewServeMux()
	mux.HandleFunc("/", base)
	mux.HandleFunc("/echo", echo)
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", host, port),
		Handler: mux,
	}

	log.Printf("Starting HTTP Server. Listening at %q", server.Addr)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		return err
		log.Printf("%v", err)
	} else {
		log.Println("Server closed!")
	}
	return nil
}
func base(w http.ResponseWriter, r *http.Request) {
	statusCode := http.StatusOK
	w.WriteHeader(statusCode)
	w.Header().Set("content-type", "text/plain; charset=utf-8")
	statusText := fmt.Sprint("GET /echo?text=hello")
	w.Write([]byte(statusText))
}

func echo(w http.ResponseWriter, r *http.Request) {
	echoResponse := EchoResponse{}
	echoResponse.AccessToken = r.Header.Get("Authorization")

	query := r.URL.Query()
	text, present := query["text"]
	if !present || len(text) == 0 {
		log.Println("text not present")
	}
	echoResponse.QueryParam = query.Get("text")

	var dst interface{}
	decode := json.NewDecoder(r.Body)
	decode.DisallowUnknownFields()
	err := decode.Decode(&dst)
	if err != nil {
		fmt.Println(err)
		err = decode.Decode(&struct{}{})
		if err != io.EOF {
			log.Printf("Request body must only contain a single JSON object", err)
			w.WriteHeader(http.StatusBadRequest)
			return

		}
	}
	if err == nil {
		echoResponse.Body = dst
	}

	w.Header().Set("content-type", "application/json")
	statusCode := http.StatusOK
	w.WriteHeader(statusCode)
	err = json.NewEncoder(w).Encode(&echoResponse)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
