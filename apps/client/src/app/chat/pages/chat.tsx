import { Chat } from "../components/chat";
import { Sidebar } from "../components/sidebar";

export const ChatPage = () => {
  return (
    <div className="h-[100vh] grid grid-cols-[350px_1fr] max-w-[1600px] mx-auto">
      <Sidebar />

      {/* <div className="h-[100vh] flex flex-col items-center justify-center px-8">
        <img src="/assets/chat-main-vector.svg" className="w-[250px]" />
        <h1 className="text-3xl font-light">Chat message system</h1>
        <p className="font-light mt-5 max-w-[400px] text-center text-muted-foreground">
          Send and receive messages end to end encrypted. No app installation
          needed.
        </p>
      </div> */}
      <Chat />
    </div>
  );
};
