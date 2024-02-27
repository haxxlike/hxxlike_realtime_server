FROM heroiclabs/nakama:3.20.1-arm
# COPY presence.yml /nakama/data/
# COPY match_1.yml /nakama/data/
# COPY match_2.yml /nakama/data/
ADD ./lua_runtime/*.yml /nakama/data/
ADD ./lua_runtime/*.lua /nakama/data/modules/
# COPY json/*.json /nakama/data/modules/
