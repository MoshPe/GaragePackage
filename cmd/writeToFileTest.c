#include <stdio.h>
#include <stdlib.h>
#include "fileWrite.h"

void printToLog(int id, char* time, char* msg)
{
    FILE *fp = fopen("logs.txt","a");
    fprintf(fp,"\n car : %d time : %s -> : %s",id, time, msg);
    fclose(fp);
}

void printDayCountToLog(int dayCount){
    FILE *fp = fopen("logs.txt","a");
    fprintf(fp,"\n\n---------- Day %d :) ----------\n\n",dayCount);
    fclose(fp);

}

void deleteFileIfExist(){
    FILE *fp = fopen("logs.txt","r");
    if (fp != NULL){
        remove("logs.txt");
    }
}

