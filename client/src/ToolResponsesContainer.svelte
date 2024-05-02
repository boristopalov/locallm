<script lang="ts">
  import Card from "./lib/Card.svelte";
  import { TopResultsEventTypes, type CardData } from "./lib/utils/types";
  export let cards: CardData[] = [];
</script>

<div class="bg-gray-600 p-2 overflow-y-auto rounded-sm h-full">
  <h3 class="text-white ml-2">Sources</h3>

  <hr class="border-gray-300 rounded-full ml-2 mb-2" />
  {#each cards as card, i}
    <Card
      data={card.data.text}
      source={card.event == TopResultsEventTypes.QueryVectorDBResult
        ? "VectorDB"
        : card.data.url}
      iconSrc={card.event == TopResultsEventTypes.QueryVectorDBResult
        ? "/images/database.svg"
        : "/images/web.svg"}
    />
  {/each}
</div>

<style>
  * {
    scrollbar-width: none;
    scrollbar-color: rgb(75 85 99) rgb(17 24 39);
    overflow-y: auto;
  }
</style>
