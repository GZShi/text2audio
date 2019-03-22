<template>
  <div id="app">
    <!-- <img src="./assets/logo.png"> -->
    <Text2Audio/>
  </div>
</template>

<script>
import HelloWorld from './components/HelloWorld'
import Text2Audio from './components/Text2Audio'

export default {
  name: 'App',
  components: {
    HelloWorld,
    Text2Audio
  }
}
</script>

<style>
* {
  outline: none;
}
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
}
</style>
