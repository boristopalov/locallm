export type WebResult = {
  name: string;
  url: string;
  text: string;
};

export type ClientQuery = {
  prompt: string;
  model: string;
  keep_alive: string;
};

export type SSEEvent = {
  event: string;
  data: string;
};

export enum EventTypes {
  WebSearch = "WebSearch",
  QueryVectorDB = "QueryVectorDB",
  Answer = "Answer",
  Error = "Error",
  Question = "Question",
}

export enum TopResultsEventTypes {
  WebSearchResult = "WebSearchResult",
  QueryVectorDBResult = "QueryVectorDBResult",
}

export type CardData = {
  data: WebResult[];
  event: string;
};
