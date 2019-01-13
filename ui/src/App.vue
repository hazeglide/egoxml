<template>
  <div>
    <Navigation></Navigation>
    <b-container fluid>
      <router-view/>
    </b-container>
  </div>
</template>

<script>
import Navigation from '@/components/Navigation.vue';
export default {
  components: { Navigation }
};
</script>

<style>
</style>
