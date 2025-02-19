import { createGroq } from "@ai-sdk/groq"
import { generateText } from "ai"

// Initialize the Groq client with the environment variable
const groq = createGroq({
  apiKey: process.env.GROQ_API_KEY!,
})

export async function POST(req: Request) {
  try {
    const { category } = await req.json()

    if (!category) {
      return new Response(JSON.stringify({ error: "Category is required" }), {
        status: 400,
        headers: { "Content-Type": "application/json" },
      })
    }

    const { text } = await generateText({
      model: groq("llama-3.3-70b-versatile"),
      prompt: `Generate a structured JSON output for Japanese vocabulary related to the theme "${category}". The output should be an array of vocabulary items, each containing kanji, romaji, english translation, and parts (individual kanji/kana components). Ensure the output is valid JSON that I can directly destructure and follows this structure:
      {
        "vocabulary": [
          {
            "kanji": "食べる",
            "romaji": "taberu",
            "english": "eat",
            "parts": [
              {
                "kanji": "食",
                "romaji": "ta"
              },
              {
                "kanji": "べる",
                "romaji": "beru"
              }
            ]
          }
        ]
      }
      Generate at least 5 vocabulary items related to the given theme. Please do not send any other info. You should only return me a raw json and nothing else,

      Here is an example of bad output:
      {
        "kanji": "晴れ",
        "romaji": "hare",
        "english": "sunny",
        "parts": [
          {
            "kanji": "晴",
            "romaji": ["seki", "haru"]
          },
          ...
        ]
        ...
      }
      The reason this is bad is because the parts of romaji that are shown do not represent the word. Instead of listing out seki haru, it should just say ha because that is what the kanji 晴 represents
    `})

    // Parse the generated text as JSON
    let vocabularyData

    try {
      vocabularyData = JSON.parse(text.split('```')[1])
    } catch (parseError) {
      console.error("Error parsing JSON:", parseError)
      return new Response(JSON.stringify({ error: "Failed to parse generated vocabulary" }), {
        status: 500,
        headers: { "Content-Type": "application/json" },
      })
    }

    return new Response(JSON.stringify(vocabularyData), {
      headers: { "Content-Type": "application/json" },
    })
  } catch (error) {
    console.error("Error:", error)
    const errorMessage = error instanceof Error ? error.message : "Failed to generate vocabulary"
    return new Response(JSON.stringify({ error: errorMessage }), {
      status: 500,
      headers: { "Content-Type": "application/json" },
    })
  }
}

