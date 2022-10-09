export function fetchData(url: string, options = {}) {
  return new Promise<Array<never>>(function (resolve, reject) {
    fetch(url, options)
      .then(async (response) => {
        const data = await response.json();
        if (!response.ok) {
          console.log(response);
          reject(response.statusText + ": " + (data && data.error));
          return;
        }
        resolve(data);
      })
      .catch((err) => {
        console.log("Catch: " + err);
        reject(err);
      });
  })
}
