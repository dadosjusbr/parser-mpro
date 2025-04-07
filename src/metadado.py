from coleta import coleta_pb2 as Coleta


def captura(year, month):
    metadado = Coleta.Metadados()
    metadado.acesso = Coleta.Metadados.FormaDeAcesso.NECESSITA_SIMULACAO_USUARIO
    metadado.extensao = Coleta.Metadados.Extensao.CSV
    metadado.estritamente_tabular = True
    metadado.tem_matricula = True
    metadado.tem_lotacao = True
    metadado.tem_cargo = True
    metadado.receita_base = Coleta.Metadados.OpcoesDetalhamento.DETALHADO
    metadado.despesas = Coleta.Metadados.OpcoesDetalhamento.DETALHADO
    metadado.outras_receitas = Coleta.Metadados.OpcoesDetalhamento.DETALHADO
    if (int(year) == 2024 and int(month) == 6):
        metadado.formato_consistente = False
    else:
        metadado.formato_consistente = True

    return metadado
