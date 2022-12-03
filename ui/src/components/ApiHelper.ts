export function fetchData(url: string, options = {}) {
  return new Promise<Array<never>>(function (resolve, reject) {
    const domain = process.env.VUE_APP_API_URL ? process.env.VUE_APP_API_URL : window.location.origin;
    fetch(domain + url, options)
      .then(async (response) => {
        const data = await response.json();
        if (!response.ok) {
          reject(response.statusText + ": " + (data && data.error));
          return;
        }
        resolve(data);
      })
      .catch((err) => {
        reject(err);
      });
  })
}
