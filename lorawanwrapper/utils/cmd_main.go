//go:build !cgo
// +build !cgo

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

// Command-line interface for non-CGO environments
func cmdMain() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <command> [args...]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Commands:\n")
		fmt.Fprintf(os.Stderr, "  printPHYPayload <payload> [key]\n")
		fmt.Fprintf(os.Stderr, "  getMType <payload>\n")
		fmt.Fprintf(os.Stderr, "  getMajor <payload>\n")
		fmt.Fprintf(os.Stderr, "  getCounter <payload>\n")
		fmt.Fprintf(os.Stderr, "  getDevEUI <payload>\n")
		fmt.Fprintf(os.Stderr, "  getJoinEUI <payload>\n")
		fmt.Fprintf(os.Stderr, "  getDevNonce <payload>\n")
		fmt.Fprintf(os.Stderr, "  getDevAddr <appKey> <payload>\n")
		fmt.Fprintf(os.Stderr, "  getDevAddrFromMACPayload <payload>\n")
		fmt.Fprintf(os.Stderr, "  generateSessionKeysFromJoins <joinRequest> <joinAccept> <appKey>\n")
		fmt.Fprintf(os.Stderr, "  generateValidMIC <payload> <key> [jaKey]\n")
		fmt.Fprintf(os.Stderr, "  marshalJsonToPHYPayload <json> [key] [nwkskey]\n")
		fmt.Fprintf(os.Stderr, "  testAppKeysWithJoinRequest <keys...> -- <payload> <generateKeys>\n")
		fmt.Fprintf(os.Stderr, "  testAppKeysWithJoinAccept <keys...> -- <payload> <generateKeys>\n")
		os.Exit(1)
	}

	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "printPHYPayload":
		if len(args) < 1 {
			fmt.Fprintf(os.Stderr, "Error: printPHYPayload requires at least 1 argument\n")
			os.Exit(1)
		}
		var key string
		if len(args) > 1 {
			key = args[1]
		}
		result := cmdPrintPHYPayload(args[0], key)
		fmt.Print(result)

	case "getMType":
		if len(args) < 1 {
			fmt.Fprintf(os.Stderr, "Error: getMType requires 1 argument\n")
			os.Exit(1)
		}
		result := cmdGetMType(args[0])
		fmt.Print(result)

	case "getMajor":
		if len(args) < 1 {
			fmt.Fprintf(os.Stderr, "Error: getMajor requires 1 argument\n")
			os.Exit(1)
		}
		result := cmdGetMajor(args[0])
		fmt.Print(result)

	case "getCounter":
		if len(args) < 1 {
			fmt.Fprintf(os.Stderr, "Error: getCounter requires 1 argument\n")
			os.Exit(1)
		}
		result := cmdGetCounter(args[0])
		fmt.Print(result)

	case "getDevEUI":
		if len(args) < 1 {
			fmt.Fprintf(os.Stderr, "Error: getDevEUI requires 1 argument\n")
			os.Exit(1)
		}
		result := cmdGetDevEUI(args[0])
		fmt.Print(result)

	case "getJoinEUI":
		if len(args) < 1 {
			fmt.Fprintf(os.Stderr, "Error: getJoinEUI requires 1 argument\n")
			os.Exit(1)
		}
		result := cmdGetJoinEUI(args[0])
		fmt.Print(result)

	case "getDevNonce":
		if len(args) < 1 {
			fmt.Fprintf(os.Stderr, "Error: getDevNonce requires 1 argument\n")
			os.Exit(1)
		}
		result := cmdGetDevNonce(args[0])
		fmt.Print(result)

	case "getDevAddr":
		if len(args) < 2 {
			fmt.Fprintf(os.Stderr, "Error: getDevAddr requires 2 arguments\n")
			os.Exit(1)
		}
		result := cmdGetDevAddr(args[0], args[1])
		fmt.Print(result)

	case "getDevAddrFromMACPayload":
		if len(args) < 1 {
			fmt.Fprintf(os.Stderr, "Error: getDevAddrFromMACPayload requires 1 argument\n")
			os.Exit(1)
		}
		result := cmdGetDevAddrFromMACPayload(args[0])
		fmt.Print(result)

	case "generateSessionKeysFromJoins":
		if len(args) < 3 {
			fmt.Fprintf(os.Stderr, "Error: generateSessionKeysFromJoins requires 3 arguments\n")
			os.Exit(1)
		}
		result := cmdGenerateSessionKeysFromJoins(args[0], args[1], args[2])
		fmt.Print(result)

	case "generateValidMIC":
		if len(args) < 2 {
			fmt.Fprintf(os.Stderr, "Error: generateValidMIC requires at least 2 arguments\n")
			os.Exit(1)
		}
		var jaKey string
		if len(args) > 2 {
			jaKey = args[2]
		}
		result := cmdGenerateValidMIC(args[0], args[1], jaKey)
		fmt.Print(result)

	case "marshalJsonToPHYPayload":
		if len(args) < 1 {
			fmt.Fprintf(os.Stderr, "Error: marshalJsonToPHYPayload requires at least 1 argument\n")
			os.Exit(1)
		}
		var key, nwkskey string
		if len(args) > 1 {
			key = args[1]
		}
		if len(args) > 2 {
			nwkskey = args[2]
		}
		result := cmdMarshalJsonToPHYPayload(args[0], key, nwkskey)
		fmt.Print(result)

	case "testAppKeysWithJoinRequest":
		separatorIndex := -1
		for i, arg := range args {
			if arg == "--" {
				separatorIndex = i
				break
			}
		}
		if separatorIndex == -1 || len(args) < separatorIndex+3 {
			fmt.Fprintf(os.Stderr, "Error: testAppKeysWithJoinRequest requires keys... -- payload generateKeys\n")
			os.Exit(1)
		}
		keys := args[:separatorIndex]
		payload := args[separatorIndex+1]
		generateKeys, err := strconv.Atoi(args[separatorIndex+2])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: generateKeys must be an integer\n")
			os.Exit(1)
		}
		result := cmdTestAppKeysWithJoinRequest(keys, payload, generateKeys)
		fmt.Print(result)

	case "testAppKeysWithJoinAccept":
		separatorIndex := -1
		for i, arg := range args {
			if arg == "--" {
				separatorIndex = i
				break
			}
		}
		if separatorIndex == -1 || len(args) < separatorIndex+3 {
			fmt.Fprintf(os.Stderr, "Error: testAppKeysWithJoinAccept requires keys... -- payload generateKeys\n")
			os.Exit(1)
		}
		keys := args[:separatorIndex]
		payload := args[separatorIndex+1]
		generateKeys, err := strconv.Atoi(args[separatorIndex+2])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: generateKeys must be an integer\n")
			os.Exit(1)
		}
		result := cmdTestAppKeysWithJoinAccept(keys, payload, generateKeys)
		fmt.Print(result)

	default:
		fmt.Fprintf(os.Stderr, "Error: Unknown command: %s\n", command)
		os.Exit(1)
	}
}

// Error handling helper
func handleError(err error, operation string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in %s: %v\n", operation, err)
		fmt.Print("Error")
		os.Exit(1)
	}
}

// JSON error response
func jsonError(message string) string {
	errorResponse := map[string]string{"error": message}
	jsonBytes, _ := json.Marshal(errorResponse)
	return string(jsonBytes)
}