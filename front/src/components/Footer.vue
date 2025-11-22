<template>
  <footer class="footer">
    <div class="container">
      <section id="footer-content" class="row text-center justify-content-center">
        <p>©{{ appName }} All Rights Reserved</p>
      </section>
    </div>
  </footer>
</template>

<script>
import common from "../common";

export default {
  data() {
    return {
      appName: common.getAppName(),
    };
  },
};
</script>

<style lang="scss">
footer {
  background-color: #c3e2ff;
  border-top: 3px dashed #9bcfff;
  text-align: center;
  padding: 15px 0;
  width: 100%;
}

#footer-content {
  margin-top: 10px;
}

@media screen and (max-width: 992px) {
  #footer-content {
    margin: 15px 0 30px;
  }
}
</style>
