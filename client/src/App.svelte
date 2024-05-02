<script lang="ts">
  import Chat from "./Chat.svelte";
  import SSERenderer from "./SSERenderer.svelte";
  import { PROMPT_API_ENDPOINT, STREAM_ENDPOINT } from "./lib/utils/vars";
  import {
    EventTypes,
    type ClientQuery,
    type SSEEvent,
    type CardData,
    TopResultsEventTypes,
    type WebResult,
  } from "./lib/utils/types";
  import ToolResponsesContainer from "./ToolResponsesContainer.svelte";
  let eventSource: EventSource | undefined;
  const lorem =
    "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum";

  // const testEvents: SSEEvent[] = [
  //   { event: EventTypes.WebSearch, data: "asdf" },
  //   { event: EventTypes.QueryVectorDB, data: "asdf123" },
  //   {
  //     event: EventTypes.Answer,
  //     data: lorem,
  //   },
  //   {
  //     event: EventTypes.Answer,
  //     data: "hey whats up",
  //   },
  //   {
  //     event: EventTypes.Error,
  //     data: "error",
  //   },
  //   {
  //     event: EventTypes.Error,
  //     data: "error",
  //   },
  //   {
  //     event: EventTypes.Answer,
  //     data: lorem,
  //   },
  //   {
  //     event: EventTypes.Answer,
  //     data: lorem,
  //   },
  // ];

  // const testWebRes: WebResult = {
  //   name: "Wikipedia",
  //   url: "wikipedia.com",
  //   text: lorem,
  // };

  // const testCards = [
  //   { data: testWebRes, event: "WebSearchResult" },
  //   { data: testWebRes, event: "WebSearchResult" },
  //   { data: testWebRes, event: "WebSearchResult" },
  //   { data: testWebRes, event: "WebSearchResult" },
  // ];
  let cards: CardData[] = [];

  let events: SSEEvent[] = [];

  function addEvent(e: SSEEvent) {
    events = [...events, e];
  }

  function addCard(c: CardData) {
    cards = [...cards, c];
  }

  const startStream = () => {
    console.log("creating event source");
    eventSource = new EventSource(STREAM_ENDPOINT);

    for (const eventType in EventTypes) {
      console.log("listening for", eventType, "event");
      eventSource.addEventListener(eventType, (e) => {
        console.log(eventType, e.data);
        addEvent({
          event: eventType,
          data: atob(e.data),
        });
      });
    }

    for (const eventType in TopResultsEventTypes) {
      console.log("listening for", eventType, "event");
      eventSource.addEventListener(eventType, (e) => {
        console.log(eventType, e.data);
        const decoded = atob(e.data);
        const webResults: WebResult[] = JSON.parse(decoded);
        console.log("parsed jSON:", JSON.parse(decoded));
        // decoded is an array of WebResult
        for (const wr of webResults) {
          addCard({
            data: wr,
            event: eventType,
          });
        }
      });
    }

    eventSource.onmessage = function (e) {
      console.log("onmessage res:", e);
    };

    eventSource.onerror = function (error) {
      console.error("SSE error:", error);
    };
  };

  const submitQuery = async (q: string) => {
    const clientQuery: ClientQuery = {
      prompt: q,
      model: "llama3-custom",
      keep_alive: "5m",
    };
    const requestOpts = {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(clientQuery),
    };
    console.log(clientQuery);
    try {
      const res = await fetch(PROMPT_API_ENDPOINT, requestOpts);
      if (!res.ok) {
        console.error(res.status);
        throw new Error("failed to make request");
      }
      if (!eventSource) {
        startStream();
      }
      addEvent({ event: EventTypes.Question, data: q });
    } catch (e) {
      console.error(e);
    }
  };
</script>

<main>
  <div id="container" class="flex bg-gray-900">
    <div class="lhs w-1/3 h-screen relative">
      <div class="pr-2 h-full">
        <ToolResponsesContainer {cards} />
      </div>
    </div>
    <div class="w-full max-h-screen relative">
      <div class="mt-32">
        <SSERenderer {events} />
        <Chat onSubmitQuery={submitQuery} />
      </div>
    </div>
  </div>
</main>
