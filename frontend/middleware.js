import { NextResponse } from "next/server";

export async function middleware(request) {
  const sessionCookie = request.cookies.get("session_id");
  const url = new URL(request.url);
  const path = url.pathname;

  const publicRoutes = ["/login", "/signup"];
  const isPublicRoute = publicRoutes.includes(path);

  // For public routes without session, allow access
  if (isPublicRoute && !sessionCookie) {
    return NextResponse.next();
  }

  // For protected routes without session, redirect to login
  if (!isPublicRoute && !sessionCookie) {
    return NextResponse.redirect(new URL("/login", request.url));
  }

  // If we have a session, verify it
  if (sessionCookie) {
    try {
      const response = await fetch("http://localhost:8080/session", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Cookie: `session_id=${sessionCookie.value}`,
        },
        credentials: "include",
      });

      const responseData = await response.json();

      // If session is valid and trying to access public route, redirect to home
      if (responseData.code === 200 && isPublicRoute) {
        return NextResponse.redirect(new URL("/", request.url));
      }

      // If session is invalid and on protected route, redirect to login
      if (responseData.code !== 200 && !isPublicRoute) {
        return NextResponse.redirect(new URL("/login", request.url));
      }
    } catch (error) {
      if (!isPublicRoute) {
        return NextResponse.redirect(new URL("/login", request.url));
      }
    }
  }

  return NextResponse.next();
}

export const config = {
  matcher: ["/((?!api|_next/static|_next/image|favicon.ico).*)"],
};
