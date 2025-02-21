export async function invokeAPI(route, body, method) {
  const options = {
    method: method,
    headers: {
      "Content-Type": "application/json",
    },
    credentials: "include",
    cache: "no-store",
  };

  if (body) {
    options.body = JSON.stringify(body);
  }

  const response = await fetch("http://localhost:8080/" + route, options);

  const data = await response.json();

  return data;
}
