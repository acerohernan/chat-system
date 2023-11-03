import { createBrowserRouter } from "react-router-dom";
import { MainLayout } from "./components/layout";
import { LoginPage } from "./app/auth/pages/login";
import { NotFoundPage } from "./pages/404";
import { ChatPage } from "./app/chat/pages/chat";

export const router = createBrowserRouter([
  {
    path: "/",
    element: <MainLayout />,
    children: [
      {
        path: "",
        element: <ChatPage />,
      },
      {
        path: "login",
        element: <LoginPage />,
      },
      {
        path: "*",
        element: <NotFoundPage />,
      },
    ],
  },
]);
