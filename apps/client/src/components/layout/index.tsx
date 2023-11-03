import { TOKEN_KEY } from "@/app/auth/constants/token";
import { useEffect } from "react";
import { Outlet, useNavigate, useSearchParams } from "react-router-dom";

export const MainLayout = () => {
  const navigate = useNavigate();
  const [params] = useSearchParams();

  const tokenParam = params.get(TOKEN_KEY);
  const storedToken = localStorage.getItem(TOKEN_KEY);

  useEffect(() => {
    if (!tokenParam && !storedToken) return navigate("/login");

    if (tokenParam) {
      localStorage.setItem(TOKEN_KEY, tokenParam);
    }

    navigate("/");
  }, [navigate, storedToken, tokenParam]);

  return <Outlet />;
};
