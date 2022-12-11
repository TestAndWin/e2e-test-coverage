export function fetchData(url: string, options: Record<string, unknown> = {}) {
  return new Promise<Array<never>>(function (resolve, reject) {
    const domain = process.env.VUE_APP_API_URL
      ? process.env.VUE_APP_API_URL
      : window.location.origin;

    const h = new Headers();
    h.append('Content-Type', 'application/json');
    h.append('Accept', 'application/json');
    const token = sessionStorage.getItem("token")
    if (token) {
      h.append('Authorization', `Bearer ${JSON.parse(token).token}`);
    }
    options["headers"] = h;

    fetch(domain + url, options)
      .then(async (response) => {
        const data = await response.json();
        if (!response.ok) {
          reject(response.statusText + ': ' + (data && data.error));
          return;
        }
        resolve(data);
      })
      .catch((err) => {
        reject(err);
      });
  });
}
