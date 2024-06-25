import { FancyAnsi } from "fancy-ansi";

async function run<T>(promise: Promise<T>): Promise<[Awaited<T> | null, any]> {
  try {
    const res = await promise;
    return [res, null];
  } catch (error) {
    return [null, error];
  }
}

interface Data {
  address: string;
  description: string;
  maxPlayers: number;
  players: number;
  playerList: string[] | null;
  version: string;
}

const url = import.meta.env.PROD ? "/data" : "http://localhost:2555/data";

const setAppHtml = (html: string) => {
  const app = document.getElementById("app");
  if (app) {
    app.innerHTML = html;
  }
};

async function App() {
  setAppHtml(`
    <p>Loading...</p>
  `);

  const [res, err] = await run(fetch(url));

  let data: Data | null = null;

  if (err) {
    setAppHtml(`
      <p>Could not get data, maybe this server is offline?</p>
    `);
  } else {
    data = await res!.json();
  }

  if (!data) {
    return;
  }

  const ansi = new FancyAnsi();
  const descriptionHtml = ansi.toHtml(data.description);

  setAppHtml(`
    <p>Address: ${data.address}</h1>
    <p>Description: ${descriptionHtml}</p>
    <p>Version: ${data.version}</p>
    <p>Players: ${data.players.toLocaleString()} / ${data.maxPlayers.toLocaleString()}</p>
    <p>Player List: ${
      data.playerList ? data.playerList.join(", ") : "Not available"
    }</p>
  `);
}

App();
