<template>
  <div class="home-layout decoration-yellow-600 p-4">
    <div class="content-container flex flex-col-reverse md:flex-row justify-between items-center gap-8 p-6">
      
      <div class="md:w-[300px] w-[180px] min-w-[90px] md:order-1 order-2">
        <img src="/images/photo_7_2025-02-19_17-13-16.jpg" alt="Image">
      </div>

      <div class="content-container description-text md:order-2 order-1">
        Lorem ipsum dolor sit amet consectetur adipisicing elit. Eligendi, ex praesentium, doloremque natus distinctio dolorem architecto, impedit corporis omnis id dolorum reiciendis ipsam perspiciatis at possimus commodi! Aspernatur, nulla corporis.
      </div>

    </div>
  </div>
</template>

<script setup>
// Removed unnecessary Header import
</script>

<style lang="scss" scoped>
.home-layout {
  background-image: url('@/assets/background.jpg');
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
  background-attachment: fixed;
  padding: 2rem;
  height: 50rem;
  width: 100%;
}

.content-container {
  max-width: 400px;
  background-color: rgba(79, 78, 78, 0.7);
  padding: 20px;
  border-radius: 8px;
}

.description-text {
  font-size: 1.125rem;
  line-height: 1.6;
  color: #333;
}
</style>
