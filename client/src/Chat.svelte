<script lang="ts">
  let inputValue = "";
  let loading = false;
  let responseMessage = "";
  export let onSubmitQuery: (q: string) => void;

  const handleKeyPress = (e: KeyboardEvent) => {
    if (
      inputValue.trim() !== "" &&
      e.key === "Enter" &&
      !e.shiftKey &&
      !e.altKey &&
      !e.ctrlKey &&
      !e.metaKey
    ) {
      e.preventDefault();
      onSubmitQuery(inputValue);
      inputValue = "";
    }
  };

  const handleClick = (e: Event) => {
    if (inputValue.trim() !== "") {
      e.preventDefault();
      onSubmitQuery(inputValue);
      inputValue = "";
    }
  };
</script>

<div class="absolute bottom-10 mx-auto left-0 right-0">
  <div class="relative flex items-center w-[500px] mx-auto">
    <button
      class="absolute right-1 p-2 bg-slate-800 rounded-full"
      on:click={handleClick}
    >
      <svg
        xmlns="http://www.w3.org/2000/svg"
        class="h-3 w-3 text-gray-400"
        fill="none"
        viewBox="0 0 24 24"
        stroke="currentColor"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
        />
      </svg>
    </button>
    <input
      type="text"
      bind:value={inputValue}
      on:keydown={handleKeyPress}
      placeholder="Ask a question"
      class="h-full w-full outline-none text-sm text-gray-700 bg-slate-300 p-2 rounded-full shadow-md"
    />
  </div>
</div>

<style>
</style>
