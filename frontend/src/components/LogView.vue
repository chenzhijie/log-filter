<template>
  <div class="log-view">
    <div class="node-btns">
        <div class="list">
          <div class="btn" :style="{ color: s == ip ? 'white' : 'black', 'background-color': s == ip ? '#2C98B6' : '#FFFFFF' }" v-for = "s in servers" @click="clickServers(s)" :key ="s">
            {{ s }}
          </div>
        </div>
        <div class="list">
            <div class="btn" :style="{ color: n == nodeIndex ? 'white' : 'black', 'background-color': n == nodeIndex ? '#2C98B6' : '#FFFFFF' }" v-for = "n in nodes" @click="clickNodes(n)" :key="n">
              {{ n }}
            </div>
        </div>
        <div class="list">
            <div class="btn" :style="{ color: t == type ? 'white' : 'black', 'background-color': t == type ? '#2C98B6' : '#FFFFFF' }" v-for = "t in types" @click="clickTypes(t)" :key="t">
              {{ t }}
            </div>
        </div>
    </div>
    <div class="log">
        <div v-html="msg"></div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Log',
  data () {
    return {
      servers: [201, 202, 203, 204, 205, 206, 207, 208, 221, 222, 223, 224, 225, 226, 227, 228],
      nodes: [1, 2, 3, 4, 5, 6, 7, 8],
      types: ['ALL', 'INFO', 'WARN', 'DEBUG', 'ERROR'],
      ip: '201',
      nodeIndex: '1',
      type: 'ALL',
      more: 0,
      msg: ''
    }
  },
  methods: {
    async clickServers (ip) {
      this.ip = ip
      await this.updateLog()
    },
    async clickNodes (n) {
      this.nodeIndex = n
      await this.updateLog()
    },
    async clickTypes (t) {
      this.type = t
      await this.updateLog()
    },
    async updateLog () {
      const result = await this.$http.get(`http://10.0.1.${this.ip}:3001/v0/log?node=${this.nodeIndex}&type=${this.type}&more=${this.more}`)
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
  },
  async created () {
    this.nodeIndex = this.$route.query.node ? this.$route.query.node : '1'
    this.type = this.$route.query.type ? this.$route.query.type : 'ALL'
    this.more = this.$route.query.more ? this.$route.query.more : '0'
    this.ip = this.$route.query.ip ? this.$route.query.ip : '201'
    await this.updateLog()
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.log-view {
  height: 100vh;
}

.node-btns {
  height: 12vh;
}

.list {
  display: flex;
  margin-bottom: 1vh;
}

.btn {
  width: 4vw;
  text-align: center;
  margin-left: 12px;
  border-radius: 4px;
  border: 1px solid;
  border-color: rgb(44, 154,182);
  color: #fff;
  background-color: rgb(44, 154,182);
  cursor: pointer;
}

.log {
  padding: 12px;
  font-size: 11pt;
  height: 88vh;
  overflow: scroll;
  /* font-weight: bold; */
  white-space: pre;
  /* color: #777777; */
  color: #F2F2F2;
  background-color: black;
}
</style>
