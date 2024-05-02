<script lang="ts">
  import { EventTypes, type SSEEvent } from "./utils/types";
  export let event: SSEEvent;
  const webSvg = "/images/web.svg";
  const dbSvg = "/images/database.svg";
  const errorSvg = "/images/error.svg";
</script>

<div class="p-2 text-xs/[1.5em] text-white mb-6">
  {#if event.event === EventTypes.WebSearch}
    <div class="flex items-center">
      <img src={webSvg} alt="web search SVG" class="size-5 mr-3" />
      <p>{event.data}</p>
    </div>
  {/if}
  {#if event.event === EventTypes.QueryVectorDB}
    <div class="flex items-center">
      <img src={dbSvg} alt="Vector DB SVG" class="size-5 mr-3" />
      <p>{event.data}</p>
    </div>
  {/if}
  {#if event.event === EventTypes.Answer}
    <div
      class="container bg-gray-600 p-3 mx-auto text-xs/[1.5em] text-white w-full shadow-lg rounded-sm"
    >
      <p>{event.data}</p>
    </div>
    <hr class="border-gray-600 rounded-full mt-5" />
  {/if}
  {#if event.event === EventTypes.Error}
    <div class="flex items-center">
      <img src={errorSvg} alt="Error SVG" class="size-5 mr-3" />
      <p>{event.data}</p>
    </div>
  {/if}

  {#if event.event === EventTypes.Question}
    <h2 class="font-semibold text-base">{event.data}</h2>
  {/if}
</div>
