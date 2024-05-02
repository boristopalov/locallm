<script lang="ts">
  export let data = "";
  export let source = "";
  export let iconSrc = "";
</script>

<div
  class="bg-gray-800/50 rounded-lg shadow-inner drop-shadow p-3 text-white mb-4"
>
  <div class="flex items-center mb-4">
    <p class="text-xs/[1.35em]">{data}</p>
  </div>

  <div class="flex text-wrap">
    <img src={iconSrc} alt="Icon" class="h-4 w-4 mr-2" />
    {#if source === "VectorDB"}
      <p class="text-xs italic">{source}</p>
    {:else}
      <a
        href={source}
        target="_blank"
        rel="noopener noreferrer"
        class="text-xs italic break-all"
      >
        {source}
      </a>
    {/if}
  </div>
</div>
