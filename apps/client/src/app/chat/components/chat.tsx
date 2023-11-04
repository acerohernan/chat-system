import { BiDotsVerticalRounded } from "react-icons/bi";
import { IoMdSend } from "react-icons/io";

export const Chat = () => {
  return (
    <div className="h-[100vh] w-full flex flex-col">
      <div className="flex items-center justify-between px-4 py-2 bg-[#F0F2F5]">
        <div className="flex gap-3">
          <div className="w-[40px] h-[40px] rounded-full bg-slate-300" />
          <div>
            <span className="block">Some</span>
            <span className="block font-light text-xs text-slate-700">
              typing..
            </span>
          </div>
        </div>
        <button className="p-2 text-xl rounded-full hover:bg-slate-200 transition-all">
          <BiDotsVerticalRounded />
        </button>
      </div>
      <div
        className="bg-[#efeae2] flex-1"
        style={{
          backgroundImage: "url('/assets/bg-chat.png')",
        }}
      ></div>
      <div className="p-4 py-3 flex gap-2">
        <input
          className="py-3 px-4 rounded-md font-light text-sm flex-1"
          placeholder="Type a message"
        />
        <button className="text-2xl px-2 rounded-full text-[#54656f] hover:bg-slate-200 transition-all">
          <IoMdSend />
        </button>
      </div>
    </div>
  );
};
