return {
  "aquasecurity/vim-trivy",
  event = "VeryLazy",
  keys = {
    { "<leader>xs", "<cmd>Trivy<cr>", desc = "Security Scan (Trivy)" },
  },
}
