package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Request struct {
	Jsonrpc string          `json:"jsonrpc"`
	Id      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type Response struct {
	Jsonrpc string          `json:"jsonrpc"`
	Id      json.RawMessage `json:"id"`
	Result  interface{}     `json:"result,omitempty"`
	Error   *Error          `json:"error,omitempty"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type InitializeResult struct {
	ProtocolVersion string `json:"protocolVersion"`
	Capabilities    struct {
		Tools struct{} `json:"tools"`
	} `json:"capabilities"`
	ServerInfo struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	} `json:"serverInfo"`
}

type Tool struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	InputSchema struct {
		Type       string `json:"type"`
		Properties map[string]struct {
			Type        string `json:"type"`
			Description string `json:"description"`
		} `json:"properties"`
		Required []string `json:"required"`
	} `json:"inputSchema"`
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	// Max capacity for scanner buffer to handle large messages
	buf := make([]byte, 1024*1024)
	scanner.Buffer(buf, 10*1024*1024)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}

		var req Request
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			sendError(nil, -32700, "Parse error")
			continue
		}

		handleRequest(&req)
	}
	
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "reading standard input: %v\n", err)
	}
}

func handleRequest(req *Request) {
	switch req.Method {
	case "initialize":
		sendResult(req.Id, InitializeResult{
			ProtocolVersion: "2024-11-05",
			ServerInfo: struct {
				Name    string `json:"name"`
				Version string `json:"version"`
			}{Name: "bash-mcp", Version: "1.0.0"},
			Capabilities: struct {
				Tools struct{} `json:"tools"`
			}{Tools: struct{}{}},
		})
	case "tools/list":
		sendResult(req.Id, map[string]interface{}{
			"tools": []Tool{
				{
					Name:        "run_bash_command",
					Description: "Execute a bash command on the host machine.",
					InputSchema: struct {
						Type       string `json:"type"`
						Properties map[string]struct {
							Type        string `json:"type"`
							Description string `json:"description"`
						} `json:"properties"`
						Required []string `json:"required"`
					}{
						Type: "object",
						Properties: map[string]struct {
							Type        string `json:"type"`
							Description string `json:"description"`
						}{
							"command": {
								Type:        "string",
								Description: "The bash command to run.",
							},
						},
						Required: []string{"command"},
					},
				},
			},
		})
	case "tools/call":
		var params struct {
			Name      string `json:"name"`
			Arguments struct {
				Command string `json:"command"`
			} `json:"arguments"`
		}
		if err := json.Unmarshal(req.Params, &params); err != nil {
			sendError(req.Id, -32602, "Invalid params")
			return
		}

		if params.Name == "run_bash_command" {
			cmd := exec.Command("bash", "-c", params.Arguments.Command)
			out, err := cmd.CombinedOutput()
			
			isError := false
			errMsg := ""
			if err != nil {
				isError = true
				errMsg = "\nError: " + err.Error()
			}

			sendResult(req.Id, map[string]interface{}{
				"content": []map[string]interface{}{
					{
						"type": "text",
						"text": string(out) + errMsg,
					},
				},
				"isError": isError,
			})
		} else {
			sendError(req.Id, -32601, "Tool not found")
		}
	default:
        // Just return an empty response for unhandled things
        if len(req.Id) > 0 {
		    sendResult(req.Id, map[string]interface{}{})
        }
	}
}

func sendResult(id json.RawMessage, result interface{}) {
	if len(id) == 0 {
		return
	}
	resp := Response{
		Jsonrpc: "2.0",
		Id:      id,
		Result:  result,
	}
	b, _ := json.Marshal(resp)
	fmt.Println(string(b))
}

func sendError(id json.RawMessage, code int, message string) {
	resp := Response{
		Jsonrpc: "2.0",
		Id:      id,
		Error: &Error{
			Code:    code,
			Message: message,
		},
	}
	b, _ := json.Marshal(resp)
	fmt.Println(string(b))
}
