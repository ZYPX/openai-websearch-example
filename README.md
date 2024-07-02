# openai-websearch-example-go

### Overview
This is a proof of concept example of tool calling following openai spec. The tool provided to the llm is called
searchWeb which allows for real-time, current information access. This example uses a [tls library](https://github.com/bogdanfinn/tls-client)
to help with scraping certain webpages. You can easily add proxy support if needed.

### Note
This code currently supports OpenRouter's api, but you can change the endpoint to OpenAI if needed in `ai.go`.

```sh
git clone <repository_url>
cd <repository_directory>
```

### Set up the environment:

Ensure you have Go installed on your machine.

### Set your OpenRouter API Key:

Replace "YOUR OPENROUTER KEY HERE" in main() with an api key from: https://openrouter.ai/.

### Usage
Run the program:
```sh
go run main.go
```