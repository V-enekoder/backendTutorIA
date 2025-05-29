package tutor

import(
	"context"
    "fmt"
    "log"
	"os"

    "google.golang.org/genai"
)

func Communication(){
	ctx := context.Background()
    client, err := genai.NewClient(ctx, &genai.ClientConfig{
        APIKey:  os.Getenv("GEMINI_API_KEY"),
        Backend: genai.BackendGeminiAPI,
    })
    if err != nil {
        log.Fatal(err)
    }

    result, err := client.Models.GenerateContent(
        ctx,
        "gemini-2.0-flash",
        genai.Text("Dame un resumen de la carrera y trayectoria de hope sandoval"),
        nil,
    )
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(result.Text())
}