/**
 * Function to invoke the API (backend) with the given route, body, method, and content type
 * Includes request caching, retries, and improved error handling
 */

// Simple in-memory cache for GET requests
const apiCache = new Map();
const CACHE_DURATION = 30 * 1000; // 30 seconds cache duration
const MAX_RETRY_COUNT = 2;
const API_BASE_URL = "http://localhost:8080/";

/**
 * Invoke API with caching, retries and improved error handling
 * @param {string} route - API route to call
 * @param {object|FormData} body - Request body
 * @param {string} method - HTTP method (GET, POST, etc)
 * @param {string} contentType - Content type header
 * @param {object} options - Additional options
 * @param {boolean} options.useCache - Whether to use cache for GET requests
 * @param {boolean} options.forceRefresh - Force refresh cache
 * @param {number} options.retryCount - Number of retries on failure
 */
export async function invokeAPI(route, body, method = 'GET', contentType, options = {}) {
  // Default options
  const {
    useCache = true,
    forceRefresh = false,
    retryCount = 0,
    timeout = 10000, // 10 second timeout by default
    headers = {}
  } = options;
  
  // Basic validation
  if (!route) {
    console.error('Invalid route provided to invokeAPI');
    return { status: 'Failed', error_msg: 'Invalid route' };
  }

  // Check cache for GET requests if caching is enabled
  const cacheKey = `${method}:${route}:${JSON.stringify(body || {})}`;
  if (method === 'GET' && useCache && !forceRefresh) {
    const cachedResponse = apiCache.get(cacheKey);
    if (cachedResponse && Date.now() < cachedResponse.expiry) {
      console.log(`Using cached response for ${route}`);
      return cachedResponse.data;
    }
  }

  // Prepare request options
  const requestOptions = {
    method: method,
    credentials: "include", // This ensures cookies are sent with the request
    cache: "no-store",
    headers: {
      // Default headers
      'Accept': 'application/json',
      ...headers // Merge custom headers
    }
  };
  
  console.log(`Request headers for ${route}:`, requestOptions.headers);

  // Add timeout with AbortController
  const controller = new AbortController();
  const timeoutId = setTimeout(() => controller.abort(), timeout);
  requestOptions.signal = controller.signal;

  // Handle body and content type
  if (body instanceof FormData) {
    // Content type for form data will be set by browser
    requestOptions.body = body;
  } else if (body) {
    requestOptions.headers["Content-Type"] = contentType || "application/json";
    if (method !== "GET") {
      requestOptions.body = JSON.stringify(body);
    }
  }
  
  try {
    const startTime = performance.now();
    console.log(`Invoking API: ${route}`);
    
    // Add a timestamp to prevent caching by the browser
    const url = new URL(API_BASE_URL + route);
    if (method === 'GET') {
      url.searchParams.append('_t', Date.now());
    }
    
    const response = await fetch(url.toString(), requestOptions);
    clearTimeout(timeoutId); // Clear the timeout
    
    const endTime = performance.now();
    console.log(`API call to ${route} took ${(endTime - startTime).toFixed(2)}ms`);
    
    let responseData;
    
    // Handle different response types
    const contentType = response.headers.get('content-type');
    
    if (contentType && contentType.includes('application/json')) {
      // JSON response
      try {
        responseData = await response.json();
      } catch (jsonError) {
        console.error(`JSON parse error: ${jsonError.message}`);
        responseData = { 
          status: 'Failed', 
          error_msg: `Failed to parse JSON response: ${jsonError.message}` 
        };
      }
    } else {
      // Non-JSON response
      try {
        const text = await response.text();
        responseData = { 
          status: response.ok ? 'Success' : 'Failed',
          data: text,
          error_msg: response.ok ? null : text
        };
      } catch (textError) {
        console.error(`Text parse error: ${textError.message}`);
        responseData = { 
          status: 'Failed', 
          error_msg: `Failed to read response: ${textError.message}` 
        };
      }
    }
    
    // Ensure consistent response format
    if (!responseData.status) {
      responseData.status = response.ok ? 'Success' : 'Failed';
    }
    
    // If response is not ok, ensure status is Failed
    if (!response.ok && responseData.status !== 'Failed') {
      responseData.status = 'Failed';
      responseData.error_msg = responseData.error_msg || `HTTP error ${response.status}`;
      
      // Handle authentication errors specifically
      if (response.status === 401) {
        responseData.code = 401;
        responseData.error_msg = 'Session expired. Please log in again.';
        responseData.unauthorized = true;
        
        // Dispatch a custom event that can be listened to by components
        if (typeof window !== 'undefined') {
          window.dispatchEvent(new CustomEvent('session-expired', { detail: responseData }));
        }
      }
    }
    
    // Cache successful GET responses
    if (method === 'GET' && useCache && response.ok) {
      apiCache.set(cacheKey, {
        data: responseData,
        expiry: Date.now() + CACHE_DURATION
      });
    }
    
    return responseData;
  } catch (error) {
    clearTimeout(timeoutId); // Clear the timeout
    
    // Handle abort errors (timeouts)
    if (error.name === 'AbortError') {
      console.error(`Request timeout for ${route} after ${timeout}ms`);
      return {
        status: "Failed",
        code: 408,
        error_msg: "Request timed out. Please try again.",
      };
    }
    
    console.error(`Network error for ${route}:`, error);
    
    // Retry logic for transient errors
    if (retryCount < MAX_RETRY_COUNT) {
      console.log(`Retrying request to ${route} (attempt ${retryCount + 1} of ${MAX_RETRY_COUNT})`);
      return invokeAPI(route, body, method, contentType, {
        ...options,
        retryCount: retryCount + 1
      });
    }
    
    return {
      status: "Failed",
      code: 500,
      error_msg: error.message || "Network error occurred",
    };
  }
}

/**
 * Clear the API cache
 * @param {string} routePrefix - Optional route prefix to clear only specific routes
 */
export function clearApiCache(routePrefix) {
  if (routePrefix) {
    // Clear only cache entries matching the prefix
    for (const key of apiCache.keys()) {
      if (key.includes(routePrefix)) {
        apiCache.delete(key);
      }
    }
    console.log(`Cleared API cache for routes matching: ${routePrefix}`);
  } else {
    // Clear all cache
    apiCache.clear();
    console.log('Cleared all API cache');
  }
}
