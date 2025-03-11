//Function to invoke the API (backend) with the given route, body, method, and content type

export async function invokeAPI(route, body, method, contentType, queryParams) {
  const options = {
    method: method,
    credentials: "include",
    cache: "no-store",
  };

  // Handle query parameters
  let url = `http://localhost:8080/${route}`;
  if (queryParams && Object.keys(queryParams).length > 0) {
    const searchParams = new URLSearchParams();
    Object.entries(queryParams).forEach(([key, value]) => {
      if (value !== null && value !== undefined) {
        searchParams.append(key, value);
      }
    });
    url += `?${searchParams.toString()}`;
  }

  if (body instanceof FormData) {
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
    const response = await fetch(url, options);
    const text = await response.text();

    try {
      const data = JSON.parse(text);
      return data;
    } catch (parseError) {
      console.error("JSON Parse Error:", parseError);
      console.error("Raw response:", text);
      return {
        code: 500,
        error_msg: "Invalid JSON response from server",
      };
    }
  } catch (error) {
    console.error("API Error:", error);
    return {
      code: 500,
      error_msg: error.message || "An unexpected error occurred",
    };
  }
}
