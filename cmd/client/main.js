const grpc = require('@grpc/grpc-js');
const protoLoader = require('@grpc/proto-loader');
const readlineSync = require('readline-sync');

// Protocol Buffers のロード
const PROTO_PATH = './path/to/openai.proto'; // プロトコルバッファのファイルパスを正しく指定
const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
  keepCase: true,
  longs: String,
  enums: String,
  defaults: true,
  oneofs: true
});
const openaiProto = grpc.loadPackageDefinition(packageDefinition).openai;

if (!openaiProto) {
  console.error("Failed to load openaiProto.");
  process.exit(1);
}

const client = new openaiProto.OpenAIService('localhost:8080', grpc.credentials.createInsecure());

function promptUser() {
  console.log("1: send Request");
  console.log("2: exit");
  const choice = readlineSync.question('please enter > ');

  if (choice === '1') {
    const inputPrompt = readlineSync.question('Please enter prompt: ');

    client.CreateChatCompletion({ prompt: inputPrompt }, (error, response) => {
      if (!error) {
        response.choices.forEach(choice => {
          console.log(choice.text); // Assuming 'text' is the field you want to print
        });
      } else {
        console.error('Connection failed.');
        console.error(error);
      }
      promptUser(); // 再度プロンプトを表示
    });
  } else if (choice === '2') {
    console.log("bye.");
  } else {
    console.log("Invalid option. Please enter 1 or 2.");
    promptUser(); // 再度プロンプトを表示
  }
}

function main() {
  console.log("start gRPC Client.");
  promptUser();
}

main();
