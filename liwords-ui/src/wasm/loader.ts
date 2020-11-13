import { Unrace } from '../utils/unrace';

// Good enough for now. If need to reload, just refresh the whole page.
class Loadable {
  private whichStep = 0;
  private fetchPromise?: Promise<Response>;

  constructor(readonly cacheKey: string, readonly path: string) {}

  startFetch = async () => {
    if (this.whichStep > 0) return;
    this.whichStep = 1;
    this.fetchPromise = fetch(this.path); // Do not await.
  };

  getArrayBuffer = async () => {
    if (this.whichStep > 1) return;
    this.startFetch(); // In case this is not done yet.
    this.whichStep = 2;
    const resp = await this.fetchPromise!;
    if (resp.ok) {
      return await resp.arrayBuffer();
    } else {
      throw new Error(`Unable to cache ${this.cacheKey}`);
    }
  };

  disownArrayBuffer = async () => {
    if (this.whichStep > 2) return;
    this.whichStep = 3; // Single-use. Assume only one Macondo worker.
    this.fetchPromise = undefined;
  };

  reset = async () => {
    this.whichStep = 0; // Allow reloading (useful when previous fetch failed).
    this.fetchPromise = undefined;
  };
}

const loadablesByLexicon: { [key: string]: Array<Loadable> } = {};

for (const { lexicons, cacheKey, path } of [
  {
    lexicons: ['CSW19', 'NWL18'],
    cacheKey: 'data/letterdistributions/english.csv',
    path: '/wasm/english.csv',
  },
  {
    lexicons: ['CSW19', 'NWL18'],
    cacheKey: 'data/strategy/default_english/leaves.idx',
    path: '/wasm/leaves.idx',
  },
  {
    lexicons: ['CSW19', 'NWL18'],
    cacheKey: 'data/strategy/default_english/preendgame.json',
    path: '/wasm/preendgame.json',
  },
  {
    lexicons: ['CSW19'],
    cacheKey: 'data/lexica/gaddag/CSW19.gaddag',
    path: '/wasm/CSW19.gaddag',
  },
  {
    lexicons: ['NWL18'],
    cacheKey: 'data/lexica/gaddag/NWL18.gaddag',
    path: '/wasm/NWL18.gaddag',
  },
]) {
  const loadable = new Loadable(cacheKey, path);
  for (const lexicon of lexicons) {
    if (!(lexicon in loadablesByLexicon)) {
      loadablesByLexicon[lexicon] = [];
    }
    loadablesByLexicon[lexicon].push(loadable);
  }
}

const unrace = new Unrace();

export interface Macondo {
  loadLexicon: (lexicon: string) => Promise<unknown>;
  precache: (loadable: Loadable) => Promise<unknown>;
  analyze: (jsonBoard: string) => Promise<string>;
}

let wrappedWorker: Macondo;

export const getMacondo = async (lexicon: string) =>
  unrace.run(async () => {
    // Allow these files to start loading.
    for (const loadable of loadablesByLexicon[lexicon] ?? []) {
      loadable.startFetch();
    }

    if (!wrappedWorker) {
      let pendings: {
        [key: string]: {
          promise: Promise<unknown>;
          res: (a: any) => void;
          rej: (a: any) => void;
        };
      } = {};

      const newPendingId = () => {
        while (true) {
          const d = String(performance.now());
          if (d in pendings) continue;

          let promRes: (a: any) => void;
          let promRej: (a: any) => void;
          const prom = new Promise((res, rej) => {
            promRes = res;
            promRej = rej;
          });

          pendings[d] = {
            promise: prom,
            res: promRes!,
            rej: promRej!,
          };

          return d;
        }
      };

      // First-time load.
      const worker = new Worker('/wasm/macondo.js');

      worker.postMessage([
        'getMacondo',
        window.RUNTIME_CONFIGURATION.macondoFilename,
      ]);

      await new Promise((res, rej) => {
        worker.onmessage = (msg) => {
          if (msg.data[0] === 'response') {
            // ["response", id, true, resp]
            // ["response", id, false] (error)
            const pending = pendings[msg.data[1]];
            if (pending) {
              if (msg.data[2]) {
                pending.res!(msg.data[3]);
              } else {
                pending.rej!(undefined);
              }
            }
          } else if (msg.data[0] === 'getMacondo') {
            // ["getMacondo", true] (ok)
            // ["getMacondo", false] (error)
            msg.data[1] ? res() : rej();
          }
        };
      });

      const sendRequest = async (req: any, transfer?: Array<Transferable>) => {
        const id = newPendingId();
        if (transfer) {
          worker.postMessage(['request', id, req], transfer);
        } else {
          worker.postMessage(['request', id, req]);
        }
        try {
          return await pendings[id].promise;
        } finally {
          delete pendings[id];
        }
      };

      class WrappedMacondo {
        loadLexicon = async (lexicon: string) => {
          return await Promise.all(
            (loadablesByLexicon[lexicon] ?? []).map((loadable) =>
              this.precache(loadable)
            )
          );
        };

        precache = async (loadable: Loadable) => {
          try {
            const arrayBuffer = (await loadable.getArrayBuffer())!;
            if (arrayBuffer) {
              await loadable.disownArrayBuffer();
              await sendRequest(
                ['precache', loadable.cacheKey, loadable.path, arrayBuffer],
                [arrayBuffer]
              );
            }
          } catch (e) {
            console.error(`failed to precache ${loadable.cacheKey}`, e);
            await loadable.reset();
            throw e;
          }
        };

        analyze = async (jsonBoard: string) => {
          return (await sendRequest(['analyze', jsonBoard])) as string;
        };
      }

      wrappedWorker = new WrappedMacondo();
    }

    await wrappedWorker.loadLexicon(lexicon);

    return wrappedWorker;
  });
