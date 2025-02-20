<template>
  <div>
    <div class="md:flex justify-between items-center space-x-6 hidden">
      <nav class="p-4 text-[#DD9735] flex space-x-6">
        <NuxtLink to="/" class="hover:underline">Explore</NuxtLink>
        <NuxtLink to="/about-us" class="hover:underline">About Us</NuxtLink>
      </nav>
      <div class="flex space-x-4">
        <button 
          class="px-2 py-1 hover:underline rounded-lg text-white bg-gradient-to-r from-[#C1842E] via-[#eec07f] to-[#5B3E16]"
          @click="goToSignUp"
        >
          SignUp
        </button>
        <button 
          class="px-2 py-1 hover:underline rounded-lg text-white bg-gradient-to-r from-[#C1842E] via-[#eec07f] to-[#5B3E16]"
          @click="goToLogin"
        >
          Login
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
// Import necessary Nuxt functions
import { useRouter } from 'vue-router';

const router = useRouter();

// Method to navigate to the SignUp page
const goToSignUp = () => {
  router.push('/signup');
};

// Method to navigate to the Login page
const goToLogin = () => {
  router.push('/login');
};
</script>

<style scoped>
/* Optional additional styling */
</style>
