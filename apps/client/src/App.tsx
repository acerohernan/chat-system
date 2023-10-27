import { useEffect } from "react";
import { Chat } from "./lib/Chat";

function App() {
  useEffect(() => {
    const url = "ws://localhost:3001/rtc";

    const chat = new Chat();

    chat.connect(url, "token");
  }, []);

  return <div>Hello world</div>;
}

export default App;
