<template>
  <div class="log">
    <div v-html="msg"></div>
  </div>
</template>

<script>
export default {
  name: 'Log',
  data () {
    return {
      msg: ''
    }
  },
  async created () {
    let nodeIndex = this.$route.query.node ? this.$route.query.node : '1'
    let type = this.$route.query.type ? this.$route.query.type : 'all'
    let more = this.$route.query.more ? this.$route.query.more : '0'
    const result = await this.$http.get(`http://localhost:3001/v0/log?node=${nodeIndex}&type=${type}&more=${more}`)
    const obj = result.body
    if (obj.result === 'success') {
      this.msg = obj.data
      this.msg = this.msg.replace(/\[0;34m/g, '<span style="color: blue">')
      this.msg = this.msg.replace(/\[0;31m/g, '<span style="color: red">')
      this.msg = this.msg.replace(/\[0;32m\[INFO \]/g, '<span style="color: green">[INFO   ]')
      this.msg = this.msg.replace(/\[0;33m/g, '<span style="color: yellow">')
      this.msg = this.msg.replace(/\[0;36m/g, '<span style="color: cyan">')
      this.msg = this.msg.replace(/\[1;35m/g, '<span style="color: pink">')
      this.msg = this.msg.replace(/\[m/g, '</span>')
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.log {
  padding: 12px;
  font-size: 11pt;
  height: 100vh;
  overflow: scroll;
  /* font-weight: bold; */
  white-space: pre;
  /* color: #777777; */
  color: #F2F2F2;
  background-color: black;
}
</style>
