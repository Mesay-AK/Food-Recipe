<template>
  <div class="w-full md:w-[280px] lg:w-[400px] min-w-[150]">
    <div class="flex justify-start items-center px-4">
      <!-- <Icon name="iconamoon:search-thin" size="20" /> -->
      <input
        type="text"
        class="py-1 pl-3 w-full rounded-md opacity-30 bg-gradient-to-r from-[#ff9d00] to-[#aa9575] placeholder-white grid-rows-2 bg-transparent outline-0 border-0 grow-1 md:pl-2 text-white"
        placeholder="Search results ..."
        v-model="searchText"
        @keyup.enter="onPressEnter"
      />
      <div class="flex justify-start items-center md:hidden lg:flex">
        <Icon
          name="solar:close-circle-bold"
          size="20"
          v-if="searchText"
          @click="clearInput"
        />
      </div>
    </div>
  </div>
</template>


<script setup>
const searchText = ref("");
const router = useRouter();

const onPressEnter = (e) => {
  e.preventDefault(); 
  if (searchText.value.trim()) {
    const q = searchText.value; 
    searchText.value = "";
    router.push("/meal?q=" + q); 
  }
};

const clearInput = () => {
  searchText.value = ""; // Clear the search text
  router.push("/meal"); // Redirect to the meals page without a query
};
</script>
