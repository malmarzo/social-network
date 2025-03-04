//Function to invoke the API (backend) with the given route, body, method, and content type

export async function invokeAPI(route, body, method, contentType) {
  const options = {
    method: method,
    credentials: "include",
    cache: "no-store",
  };

  if (body instanceof FormData) {
    // Content type for form data will be set by browser
    options.body = body;
  } else {
    options.headers = {
      "Content-Type": contentType || "application/json",
    };
    if (method !== "GET") {
      options.body = JSON.stringify(body);
    }
  }
  
  try {
    const response = await fetch("http://localhost:8080/" + route, options);
    const data = await response.json();
    return data;
  } catch (error) {
    console.error("API Error:", error);
    return {
      code: 500,
      error_msg: error.message || "An unexpected error occurred",
    };
  }
}

