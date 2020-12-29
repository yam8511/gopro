<template>
  <img style="width: 50%" alt="Vue logo" src="./assets/kubernetes-logo.png" />
  <div
    v-if="!start"
    class="w3-container w3-margin w3-larger"
    @click="
      () => {
        start = true;
      }
    "
  >
    <h1 class="w3-button w3-blue w3-round-large">Start Setup</h1>
  </div>
  <HelloWorld v-else msg="Cluster Setup" />
</template>

<script>
import HelloWorld from "./components/HelloWorld.vue";

export default {
  name: "App",
  components: {
    HelloWorld,
  },
  data() {
    return {
      start: false,
    };
  },
};
</script>
