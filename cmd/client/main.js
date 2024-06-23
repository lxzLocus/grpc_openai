const grpc = require('@grpc/grpc-js');
const protoLoader = require('@grpc/proto-loader');
const readlineSync = require('readline-sync');
const packageDefinition = protoLoader.loadSync('../../api/protocols.proto', {
  keepCase: true,
  longs: String,
  enums: String,
  defaults: true,
  oneofs: true
});

const proto = grpc.loadPackageDefinition(packageDefinition).openai;

function main() {
  try{
    const client = new proto.OpenAIService('localhost:8080', grpc.credentials.createInsecure());
    

    while (true) {
      console.log("1: send Request");
      console.log("2: exit");
      const choice = readlineSync.question('please enter > ');

      if (choice === '1') {
        const inputPrompt = readlineSync.question('Please enter prompt: ');

        client.CreateChatCompletion({ prompt: inputPrompt }, (error, response) => {
          if (!error) {
            response.choices.forEach(choice => {
              //console.log(choice.text); 
            });
          } else {
            console.error('Connection failed.');
            console.error(error);
          }
        });
      } else if (choice === '2') {
        console.log("bye.");
        break;
      } else {
        console.log("Invalid option. Please enter 1 or 2.");
      }
    }
  }catch(error){
    console.log(error);
  }
}


main();
