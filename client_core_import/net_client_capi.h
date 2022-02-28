#include <stdint.h>

typedef uintptr_t netclient_handle_t;

netclient_handle_t net_client_new(char* address);

void net_client_delete(netclient_handle_t h);

int net_client_send_login_req(netclient_handle_t h, char* account, char* password);

int net_client_send_game_enter_req(netclient_handle_t h, char* account, char* session_token);

int net_client_send_time_sync_req(netclient_handle_t h);

int net_client_send_tank_move_req(netclient_handle_t h, int dir);

int net_client_send_tank_update_pos_req(netclient_handle_t h, int move_state, int x, int y, int dir, int speed);

int net_client_send_tank_stop_move_req(netclient_handle_t h);

int net_client_send_tank_change_req(netclient_handle_t h);

int net_client_send_tank_restore_req(netclient_handle_t h);
