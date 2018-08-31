<script>
export default {
  name: 'auth',
  props: ['auth'],
  data () {
    this.auth.handleAuthentication()
    return {}
  }
}
</script>