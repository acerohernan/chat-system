import { BiDotsVerticalRounded } from "react-icons/bi";
import { BsCheckAll } from "react-icons/bs";

interface Item {
  name: string;
  lastMessage: ChatMessage;
  pendingMessages: number;
}

interface ChatMessage {
  content: string;
  timestamp: number;
}

const items: Item[] = [
  {
    name: "Someone",
    lastMessage: {
      content: "text example largemessagefsdsdf dfsdfsdfsdfdsfsdfds",
      timestamp: Date.now(),
    },
    pendingMessages: 1,
  },
  {
    name: "Someone",
    lastMessage: {
      content: "text example largemessage",
      timestamp: Date.now(),
    },
    pendingMessages: 1,
  },
  {
    name: "Someone",
    lastMessage: {
      content: "text example largemessage",
      timestamp: Date.now(),
    },
    pendingMessages: 1,
  },

  {
    name: "Someone",
    lastMessage: {
      content: "text example largemessage",
      timestamp: Date.now(),
    },
    pendingMessages: 1,
  },
  {
    name: "Someone",
    lastMessage: {
      content: "text example largemessage",
      timestamp: Date.now(),
    },
    pendingMessages: 1,
  },

  {
    name: "Someone",
    lastMessage: {
      content: "text example largemessage",
      timestamp: Date.now(),
    },
    pendingMessages: 1,
  },
  {
    name: "Someone",
    lastMessage: {
      content: "text example largemessage",
      timestamp: Date.now(),
    },
    pendingMessages: 1,
  },
];

export const Sidebar = () => {
  return (
    <div className="w-100 h-[100vh] bg-white border-[#eeeeee] border-r">
      <div className="flex items-center justify-between px-4 py-2 bg-[#F0F2F5]">
        <div className="w-[40px] h-[40px] rounded-full bg-slate-300" />
        <button className=" p-2 text-xl rounded-full hover:bg-slate-200 transition-all">
          <BiDotsVerticalRounded />
        </button>
      </div>
      <div>
        {items.map((i) => (
          <Contact item={i} />
        ))}
      </div>
    </div>
  );
};

const Contact: React.FC<{ item: Item }> = ({ item }) => {
  return (
    <div className="flex items-center py-3 px-4 gap-3 border-b border-[#f5f6f6] hover:bg-bg_prim transition-all ease-in-out cursor-pointer">
      <div className="w-[49px] h-[49px] rounded-full bg-slate-300 flex-shrink-0 " />
      <div className="overflow-hidden flex-1">
        <div className="flex items-center justify-between">
          <span className="block">{item.name}</span>
          <div className="text-xs font-light">8:37 PM</div>
        </div>
        <div className="font-light text-sm flex gap-1">
          <div className="text-blue-300 text-xl">
            <BsCheckAll />
          </div>
          <div className="line-clamp-1">{item.lastMessage.content}</div>
        </div>
      </div>
    </div>
  );
};
