<script lang="ts">
  import { onDestroy, onMount } from "svelte";
  import type { SSEEvent } from "./lib/utils/types";
  import EventComponent from "./lib/EventComponent.svelte";
  export let events: SSEEvent[] = [];
</script>

<div
  class="w-[500px] 2xl:h-[600px] md:h-[400px] max-h-full mx-auto overflow-y-auto"
>
  {#each events as event}
    <EventComponent {event} />
  {/each}
</div>

<style>
  * {
    scrollbar-width: thin;
    scrollbar-color: rgb(75 85 99) rgb(17 24 39);
    overflow-y: auto;
  }
</style>
