import { cookies } from "next/headers";
import { redirect } from "next/navigation";


export async function authFetch(
  path: string,
  init: RequestInit = {}
) {
  const token = (await cookies()).get("token")?.value;

  if (!token) {
    redirect("/login");
  }

  const apiBase =
  typeof window === "undefined"
    ? process.env.API_INTERNAL_URL
    : process.env.NEXT_PUBLIC_API_URL;

  console.log("Making request to:", `${apiBase}${path} from ${typeof window === "undefined" ? "server" : "client"}`);

  const res = await fetch(`${apiBase}${path}`, {
    ...init,
    headers: {
      ...(init.headers || {}),
      Authorization: `Bearer ${token}`,
      "Content-Type": "application/json",
    },
  });

  if (res.status === 401) {
    redirect("/login");
  }

  return res;
}
