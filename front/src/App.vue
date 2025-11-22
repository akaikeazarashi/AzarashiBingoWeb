<template>
  <navbar></navbar>

  <div id="main-area" class="container-fluid">
    <RouterView v-slot="{ Component }">
      <transition name="fade" mode="out-in">
        <component :is="Component" />
      </transition>
    </RouterView>
  </div>

  <Footer></Footer>
</template>

<script setup lang="ts">
import Navbar from "./components/Navbar.vue";
import Footer from "./components/Footer.vue";
</script>

<style lang="scss">
#app {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  color: #565656;
  background-color: #ebf6fd;
  font-family: "ヒラギノ角ゴ Pro W3", "Hiragino Kaku Gothic Pro", "メイリオ", "ＭＳ Ｐゴシック", "MS PGothic", Sans-Serif;
}

#app {
  p {
    font-size: 1rem;
  }
  h1 {
    font-size: 2rem;
  }

  h2 {
    font-size: 1.5rem;
  }

  h3 {
    font-size: 1.2rem;
  }

  .btn {
    border-radius: 50rem !important;
  }
}

#main-area {
  margin-bottom: 40px;
  flex: 1;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.gf-content-box {
  padding: 10px 20px;
  background: #e1f3ff;
  box-shadow: 0px 0px 0px 10px #e1f3ff;
  border: dashed 3px #adcdff;
  border-radius: 8px;
}

.gf-content-img {
  border-radius: 10px;
  box-shadow: 0px 0px 15px -5px #56460b;
}

@media screen and (max-width: 992px) {
  #app h1 {
    font-size: 2.5rem;
  }

  #app h2 {
    font-size: 2rem;
  }

  #app h3 {
    font-size: 1.75rem;
  }
}

@media screen and (max-width: 550px) {
  #app p {
    font-size: 1rem;
  }

  #app h1 {
    font-size: 2rem;
  }

  #app h2 {
    font-size: 1.75rem;
  }

  #app h3 {
    font-size: 1.5rem;
  }
}
</style>
