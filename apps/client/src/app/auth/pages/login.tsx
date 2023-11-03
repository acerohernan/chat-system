import { AiOutlineGoogle, AiOutlineMessage } from "react-icons/ai";

import { Button } from "@/components/ui/button";

export const LoginPage = () => {
  return (
    <div className="w-full h-[100vh] bg-slate-200 flex items-center justify-center">
      <div className="flex flex-col items-center">
        <div className="text-4xl mb-2">
          <AiOutlineMessage />
        </div>
        <h1>Message system</h1>
        <Button
          size="lg"
          type="button"
          asChild
          className="flex items-center mt-8 justify-center gap-2"
        >
          <a
            className="cursor-pointer"
            href={`${import.meta.env.VITE_API_URL}/auth/google`}
          >
            <div className="text-2xl">
              <AiOutlineGoogle />
            </div>
            Continue with Google
          </a>
        </Button>
      </div>
    </div>
  );
};
