import { Configuration, Round, RoundApi } from "../api";

const configuration = new Configuration({
  // TODO: Config this on build or use
  // basePath: window.location.origin
  basePath: "http://localhost:8037",
});

const roundAPI = new RoundApi(configuration);

export function getRoundByID(id: string): Promise<Round> {
  return roundAPI.roundRoundIDGet({
    roundID: id,
  });
}

export { roundAPI };
