# openai-websearch-example-go

### Overview
This is a proof of concept example of tool calling following openai spec. The tool provided to the llm is called
searchWeb which allows for real-time, current information access. This example uses a [tls library](https://github.com/bogdanfinn/tls-client)
to help with scraping certain webpages. Proxy support can be added.

### Note
This code currently supports OpenRouter's api, but you can change the endpoint and other parameters to OpenAI if needed
in `ai.go`.

Responses are streamed.

### Set up the environment:
```sh
git clone <repository_url>
cd <repository_directory>
```

Ensure you have Go installed on your machine.

### Set your OpenRouter API Key:

Replace "YOUR OPENROUTER KEY HERE" in main() with an api key from: https://openrouter.ai/.

### Usage
Run the program:
```sh
go run main.go
```
Type `end` to exit the program.

### License
This project is licensed under the MIT License - see the `LICENSE` file for details.
