package chat

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/V-enekoder/backendTutorIA/config"
	"github.com/V-enekoder/backendTutorIA/src/schema"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func ProcessPromptService(promptText string, uri UriParams) (string, error) {
	apiKey := config.GetGeminiAPIKey()

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Printf("Error creando el cliente de Gemini: %v", err)
		return "", errors.New("no se pudo inicializar el cliente de Gemini")
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")

	finalPrompt := promptText

	if uri.ID >= 1 && uri.ID <= uint(len(promptContexts)) {
		contextIndex := uri.ID - 1
		contextText := promptContexts[contextIndex]
		finalPrompt = fmt.Sprintf("Instrucción de contexto: %s\n\n---\n\nPregunta del usuario: %s. En la respuesta omite los saltos de línea y otros caracteres no visibles. Y los asteriscos y otros elementos de markdown.", contextText, promptText)
		log.Printf("ID %d válido. Se ha añadido el contexto: '%s'", uri.ID, contextText)
	} else {
		// Si el ID no es válido (ej: 0, -1, 21, etc.), simplemente informamos y continuamos sin el contexto.
		log.Printf("ID %d fuera de rango [1-%d]. Se usará el prompt sin contexto adicional.", uri.ID, len(promptContexts))
	}

	resp, err := model.GenerateContent(ctx, genai.Text(finalPrompt))
	if err != nil {
		log.Printf("Error generando contenido con Gemini: %v", err)
		return "", errors.New("error al comunicarse con la API de Gemini")
	}

	var responseBuilder strings.Builder
	if resp != nil && len(resp.Candidates) > 0 && resp.Candidates[0].Content != nil {
		for _, part := range resp.Candidates[0].Content.Parts {
			if txt, ok := part.(genai.Text); ok {
				responseBuilder.WriteString(string(txt))
			}
		}
	} else {
		log.Println("Respuesta de Gemini vacía o inesperada")
		return "", errors.New("respuesta vacía o inesperada de Gemini")
	}

	chat := schema.Chat{
		UserID:    uri.UserID,
		ContextID: uri.ID,
		Prompt:    promptText,
		Response:  responseBuilder.String(),
	}
	if err := CreateChatRepository(chat); err != nil {
		log.Printf("Error guardando el chat en la base de datos: %v", err)
		return "", errors.New("error al guardar el chat en la base de datos")
	}

	return responseBuilder.String(), nil
}

var promptContexts = []string{
	"Eres un asistente general para la redacción de tesis. Pregunta al usuario en qué parte específica de su investigación necesita ayuda (problema, metodología, redacción, etc.) y ofrécele orientación inicial general. Responde siempre en español.",
	// 3. Capítulo 1
	"Enfócate en el Capítulo 1 (El Problema). Ayuda al usuario a formular el planteamiento del problema, la pregunta de investigación, los objetivos (general y específicos) y la justificación de su tesis. Responde siempre en español.",
	// 4. Capítulo 2
	"Tu rol es ayudar a construir el Capítulo 2 (Marco Teórico). Asiste al usuario en la búsqueda de antecedentes, la definición de las bases teóricas, y la conceptualización de las variables de la investigación. Responde siempre en español.",
	// 5. Capítulo 3
	"Concéntrate en el Capítulo 3 (Marco Metodológico). Guía al usuario para definir el tipo y diseño de la investigación, la población y muestra, y las técnicas e instrumentos de recolección de datos. Responde siempre en español.",
	// 6. Capítulo 4
	"Asiste en la elaboración del Capítulo 4 (Análisis e Interpretación de Resultados). Ayuda al usuario a estructurar la presentación de los datos obtenidos y a realizar un análisis coherente con la metodología planteada. Responde siempre en español.",
	// 7. Capítulo 5
	"Ayuda a redactar el Capítulo 5 (Conclusiones y Recomendaciones). Asiste al usuario para sintetizar los hallazgos, contrastarlos con los objetivos, y proponer conclusiones claras y recomendaciones pertinentes. Responde siempre en español.",
	// 8. Delimitación Espacio Temporal
	"Tu tarea es ayudar a definir el alcance. Guía al usuario para establecer claramente la delimitación espacial (lugar geográfico) y temporal (período de tiempo) de su estudio de investigación. Responde siempre en español.",
	// 9. Temáticas posibles
	"Actúa como un generador de ideas. Basado en el área de estudio del usuario, proporciona una lista de posibles temáticas de investigación que sean relevantes, originales y factibles. Responde siempre en español.",
	// 10. Hechos conocidos
	"Tu rol es ayudar a establecer el estado del arte. Ayuda al usuario a listar y resumir los hechos conocidos, datos y el conocimiento ya establecido sobre su tema de investigación. Responde siempre en español.",
	// 11. Síntomas
	"Utilizando la metodología del 'árbol de problemas', ayuda al usuario a identificar y describir los 'síntomas' o efectos observables del problema de investigación que ha elegido. Responde siempre en español.",
	// 12. Posibles Causas
	"Siguiendo con la metodología del 'árbol de problemas', asiste al usuario en la lluvia de ideas y formulación de las posibles causas (directas e indirectas) que originan el problema de investigación. Responde siempre en español.",
	// 13. Consecuencias
	"En el marco del 'árbol de problemas', ayuda al usuario a articular y describir las consecuencias o efectos a mediano y largo plazo si el problema de investigación no se soluciona. Responde siempre en español.",
	// 14. Lo Investigable
	"Actúa como un asesor de viabilidad. Ayuda al usuario a transformar una idea general en un problema de investigación específico, medible, alcanzable, relevante y con un plazo definido (SMART). Responde siempre en español.",
	// 15. Referentes
	"Concéntrate en los referentes teóricos y empíricos (antecedentes). Ayuda al usuario a identificar autores clave, teorías fundamentales y estudios previos que sirvan de soporte para su investigación. Responde siempre en español.",
	// 16. Bases Legales
	"Tu tarea es asistir en la identificación de las bases legales. Ayuda al usuario a encontrar y citar las leyes, reglamentos, normativas o decretos que enmarcan y dan sustento jurídico a su tema de tesis. Responde siempre en español.",
	// 17. Título Tentativo
	"Ayuda al usuario a crear un 'título de trabajo' o tentativo. Este título no tiene que ser perfecto, pero debe capturar la esencia de lo que se quiere investigar para empezar a trabajar. Responde siempre en español.",
}

func ProcessPromptWithFileService(ctx context.Context, promptText string, fileData *GeminiFileData, uri UriParams) (string, error) {
	if fileData == nil || len(fileData.Data) == 0 {
		return "", errors.New("ProcessPromptWithFileService requiere datos de archivo válidos")
	}
	if fileData.MIMEType == "" {
		return "", errors.New("ProcessPromptWithFileService requiere un MIMEType para el archivo")
	}

	apiKey := config.GetGeminiAPIKey()

	if ctx == nil {
		log.Println("Advertencia: ProcessPromptWithFileService recibió un contexto nulo, usando context.Background().")
		ctx = context.Background()
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Printf("Error creando el cliente de Gemini: %v", err)
		return "", errors.New("no se pudo inicializar el cliente de Gemini")
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")

	finalPrompt := promptText

	if uri.ID >= 1 && uri.ID <= uint(len(promptContexts)) {
		contextIndex := uri.ID - 1
		contextText := promptContexts[contextIndex]

		finalPrompt = fmt.Sprintf("Instrucción de contexto: %s\n\n---\n\nPregunta del usuario: %s. En la respuesta omite los saltos de línea"+
			"y otros caracteres no visibles. Y los asteriscos y otros elementos de markdown.", contextText, promptText)
		log.Printf("ID %d válido. Se ha añadido el contexto: '%s'", uri.ID, contextText)
	} else {
		log.Printf("ID %d fuera de rango [1-%d]. Se usará el prompt sin contexto adicional.", uri.ID, len(promptContexts))
	}

	var parts []genai.Part
	parts = append(parts, genai.Blob{MIMEType: fileData.MIMEType, Data: fileData.Data})
	parts = append(parts, genai.Text(finalPrompt))

	resp, err := model.GenerateContent(ctx, parts...)
	if err != nil {
		log.Printf("Error generando contenido con Gemini: %v", err)
		return "", errors.New("error al comunicarse con la API de Gemini")
	}

	var responseBuilder strings.Builder
	if resp != nil && len(resp.Candidates) > 0 && resp.Candidates[0].Content != nil {
		for _, part := range resp.Candidates[0].Content.Parts {
			if txt, ok := part.(genai.Text); ok {
				responseBuilder.WriteString(string(txt))
			}
		}
	} else {
		log.Println("Respuesta de Gemini vacía o inesperada")
		return "", errors.New("respuesta vacía o inesperada de Gemini")
	}

	chat := schema.Chat{
		UserID:    uri.UserID,
		ContextID: uri.ID,
		Prompt:    promptText,
		Response:  responseBuilder.String(),
	}
	if err := CreateChatRepository(chat); err != nil {
		log.Printf("Error guardando el chat en la base de datos: %v", err)
		return "", errors.New("error al guardar el chat en la base de datos")
	}

	return responseBuilder.String(), nil
}

func determineMIMETypeService(filename string) string {
	if strings.HasSuffix(filename, ".pdf") {
		return "application/pdf"
	} else if strings.HasSuffix(filename, ".png") {
		return "image/png"
	} else if strings.HasSuffix(filename, ".jpg") || strings.HasSuffix(filename, ".jpeg") {
		return "image/jpeg"
	} else if strings.HasSuffix(filename, ".txt") {
		return "text/plain"
	}

	return ""
}

func GetChatHistoryService(userID uint, contextID uint) ([]ChatResponseDTO, error) {
	chats, err := GetChatsByUserIDAndContextIDRepository(userID, contextID)
	if err != nil {
		log.Printf("Error al obtener el historial del chat desde el repositorio: %v", err)
		return nil, err // Devuelve el error para que el controlador lo maneje.
	}
	var chatDTOs []ChatResponseDTO
	for _, chat := range chats {
		chatDTOs = append(chatDTOs, ChatResponseDTO{
			ID:        chat.ID,
			Prompt:    chat.Prompt,
			Response:  chat.Response,
			CreatedAt: chat.CreatedAt,
		})
	}

	return chatDTOs, nil
}
