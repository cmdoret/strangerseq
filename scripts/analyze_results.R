# Requires the raincloudplots and patchwork packages
# cmdoret, 20210305

# install.packages('devtools')
# devtools::install_github("thomasp85/patchwork")
# install.packages("remotes")
# remotes::install_github('jorvlan/raincloudplots')
library(raincloudplots)
library(readr)
library(dplyr)
library(patchwork)

# Scores have different orders of magnitude between parameter
# combination, so we first standardize them in each combination.
# Then we compute mean for each set separately in each combo.
df = read_tsv('out/bench_data.tsv') %>%
    group_by(k, seq_len, gc_weight, mode) %>%
    mutate(full_score=(full_score - mean(full_score)) / sd(full_score))%>%
    ungroup() %>%
    group_by(k, seq_len, gc_weight, mode, subset) %>%
    summarize(mean_score=mean(full_score))

df_sub = df %>% filter(mode == 'divergent')

df_sub <- data_1x1(
  array_1 = df_sub$mean_score[df_sub$subset=='seq'],
  array_2 = df_sub$mean_score[df_sub$subset=='control'],
  jit_distance = .09,
  jit_seed = 321)

raincloud_2 <- raincloud_1x1_repmes(
  data = df_sub,
  colors = (c('dodgerblue', 'darkorange')),
  fills = (c('dodgerblue', 'darkorange')),
  line_color = 'gray',
  line_alpha = .3,
  size = 1,
  alpha = .6,
  align_clouds = TRUE) +
  scale_x_continuous(
    breaks=c(1,2), labels=c("Control", "Optimized"), limits=c(0, 3)
  ) +
  xlab("Condition") + 
  ylab("Standardized mean score for parameter combination") +
  ggtitle("Control vs optimized divergent sequences") +
  theme_classic()

pdf('results_benchmark.pdf')
raincloud_2
dev.off()