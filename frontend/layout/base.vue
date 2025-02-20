<template>
  <div class="w-full min-h-screen flex flex-col">
    <Header /> <!-- Ensure this is included if needed -->
    <NuxtPage class="flex-grow" />

    <footer
      class="flex w-full h-auto px-24 py-6 justify-between items-center bg-green-600 mt-auto"
    >
      <div class="flex items-center justify-start text-white">
        &copy; Copyright, SDH 2023, All Rights Reserved
      </div>
      <div class="flex items-center justify-between gap-6 text-white">
        <NuxtLink to="/">
          <Icon name="ph:instagram-logo-light" size="25" />
        </NuxtLink>
        <NuxtLink to="/">
          <Icon name="bxl:twitter" size="25" />
        </NuxtLink>
        <NuxtLink to="/">
          <Icon name="ic:baseline-facebook" size="25" />
        </NuxtLink>
        <NuxtLink to="/">
          <Icon name="uil:youtube" size="25" />
        </NuxtLink>
      </div>
    </footer>
  </div>
</template>

<script setup>
import Header from "@/components/headers/Header.vue";
</script>
