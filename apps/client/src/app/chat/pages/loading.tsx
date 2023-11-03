import { AiOutlineMessage } from "react-icons/ai";

export const LoadingPage = () => {
  return (
    <div className="w-full h-[100vh] bg-slate-200 flex items-center justify-center">
      <div className="flex flex-col items-center">
        <div className="text-4xl mb-2">
          <AiOutlineMessage />
        </div>
        <h1>Message system</h1>
        <p className="mt-8">Loading chats...</p>
      </div>
    </div>
  );
};
